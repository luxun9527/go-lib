package client

import (
	"context"
	"github.com/spf13/cast"
	"go-lib/net/grpc/pb/grpcdemo"
	"go-lib/net/grpc/pb/grpcdemo/folder"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/resolver"
	"log"
	"math/rand"
	"testing"
	"time"
)

const WeightLBName = "weight_lb"

// 自定义负载均衡，加权轮询
type weightConf struct {
	addr   string
	weight int32
}

// 自定义 Picker
type weightsPicker struct {
	subConns []balancer.SubConn // 连接列表
	weights  []*weightConf      // 权重列表
	current  uint32             // 当前轮询位置
}

func (p *weightsPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	// 如果没有可用的连接，返回错误
	if len(p.subConns) == 0 {
		return balancer.PickResult{}, nil
	}

	// 计算当前要选择的连接，权重越大的连接被选择的概率越大
	totalWeight := float64(0)
	for _, w := range p.weights {
		totalWeight += float64(w.weight)
	}

	if totalWeight == 0 {
		return balancer.PickResult{}, nil
	}

	// 构建累积权重
	cumulativeWeights := make([]int32, len(p.weights))
	cumulativeWeights[0] = p.weights[0].weight
	for i := 1; i < len(p.weights); i++ {
		cumulativeWeights[i] = cumulativeWeights[i-1] + p.weights[i].weight
	}

	// 生成随机数并选择值
	rand.Seed(time.Now().UnixNano())
	randomValue := rand.Float64() * totalWeight
	idx := 0
	// 找到随机数落入的范围
	for i, cw := range cumulativeWeights {
		if randomValue <= float64(cw) {
			idx = i
			break
		}
	}
	// 计算加权轮询的下标

	md, _ := metadata.FromIncomingContext(info.Ctx)
	log.Printf("Pick subconn %d with metadata %v and weight %v total weight %v randomValue %v", idx, md, p.weights[idx], totalWeight, randomValue)

	// 返回选择的连接
	return balancer.PickResult{SubConn: p.subConns[idx]}, nil
}

// 负载均衡器构建器
type weightsPickerPickerBuilder struct {
	weightConfig map[string]int32
}

func (wp *weightsPickerPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}

	// 这里我们假设每个连接的权重是相同的，或者你可以根据实际情况传递权重
	var subConns []balancer.SubConn
	var weights []*weightConf
	for sc, addr := range info.ReadySCs {
		subConns = append(subConns, sc)
		weight := addr.Address.Attributes.Value("weight")
		if weight == nil {
			weight = int32(1)
		}
		log.Printf("Adding subconn %s with weight %v", addr.Address.Addr, weight)
		weights = append(weights, &weightConf{
			addr:   addr.Address.Addr,
			weight: cast.ToInt32(weight),
		}) // 这里给每个连接设置默认权重为1，实际应用中可以根据情况设置
	}

	return &weightsPicker{
		subConns: subConns,
		weights:  weights,
	}
}

// 自定义负载均衡
func newWeightBalancerBuilder() balancer.Builder {
	return base.NewBalancerBuilder(WeightLBName, &weightsPickerPickerBuilder{}, base.Config{HealthCheck: true})
}

func TestBalancer(t *testing.T) {
	balancer.Register(newWeightBalancerBuilder())
	resolverBuilder := &customBuilder{}
	resolver.Register(resolverBuilder)

	conn, err := grpc.Dial(
		"dns:///xxx.xxx.com",
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"weight_lb"}`),
		grpc.WithResolvers(resolverBuilder),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Panicf("dial connection failed err =%v", err)
	}
	cli := grpcdemo.NewGrpcDemoClient(conn)

	addrCount := map[string]int32{}

	for i := 0; i < 200; i++ {
		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("key", "value"))
		resp, err := cli.DemoImport(ctx, &folder.ImportedMessage{
			ImportedMessage: "test",
		})
		if err != nil {
			log.Printf("call demoImport failed err %v", err)
			continue
		}
		addrCount[resp.CustomMessage] = addrCount[resp.CustomMessage] + 1
		log.Printf("调用地址统计: %v", addrCount)
		//2024/11/28 00:30:55 调用地址统计: map[8897:113 8898:63 8899:24]
		time.Sleep(time.Millisecond * 4)
	}
}
