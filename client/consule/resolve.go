package consule

import (
	"fmt"
	"regexp"
	"sync"

	"google.golang.org/grpc/resolver"

	"github.com/hashicorp/consul/api"
)

func New() resolver.Builder {
	return &Builder{}
}

type Builder struct {

}


type Resolver struct {
	address              string
	wg                   sync.WaitGroup
	cc                   resolver.ClientConn
	name                 string
	disableServiceConfig bool
	lastIndex            uint64
}


func (cb *Builder) Scheme() string {
	return "consul"
}

func (rs *Resolver) ResolveNow(opt resolver.ResolveNowOption) {
	fmt.Println(opt)
}

func (rs *Resolver) Close() {
}

func (rs *Resolver) watcher() {
	fmt.Printf("calling consul watcher\n")
	config := api.DefaultConfig()
	config.Address = rs.address
	client, err := api.NewClient(config)
	if err != nil {
		fmt.Printf("error create consul client: %v\n", err)
		return
	}
	for {
		services, metainfo, err := client.Health().Service(rs.name, rs.name, true, &api.QueryOptions{WaitIndex: rs.lastIndex})
		if err != nil {
			fmt.Printf("error retrieving instances from Consul: %v", err)
		}

		rs.lastIndex = metainfo.LastIndex
		var newAddrs []resolver.Address
		for _, service := range services {
			addr := fmt.Sprintf("%v:%v", service.Service.Address, service.Service.Port)
			newAddrs = append(newAddrs, resolver.Address{Addr: addr})
		}
		state := resolver.State{
			Addresses:newAddrs,
			ServiceConfig:rs.name,
		}
		fmt.Printf("newAddrs: %v\n", state)
		rs.cc.UpdateState(state)
	}
}

func (cb *Builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {
	fmt.Printf("build\n")
	fmt.Printf("target: %+v\n", target)
	host, port, name, err := parseTarget(fmt.Sprintf("%s/%s", target.Authority, target.Endpoint))
	if err != nil {
		return nil, err
	}
	cr := &Resolver{
		address:              fmt.Sprintf("%s%s", host, port),
		name:                 name,
		cc:                   cc,
		disableServiceConfig: opts.DisableServiceConfig,
	}
	fmt.Printf("target: %+v\n", cr)
	cr.wg.Add(1)
	go cr.watcher()
	return cr,nil
}

func parseTarget(target string) (host, port, name string, err error) {
	fmt.Printf("target uri: %v\n", target)
	if target == "" {
		return "", "", "", fmt.Errorf("target is null")
	}

	exp,err := regexp.Compile("^([A-z0-9.]+)(:[0-9]{1,5})?/([A-z_]+)$")
	if err != nil{
		return
	}
	groups := exp.FindStringSubmatch(target)
	host = groups[1]
	port = groups[2]
	name = groups[3]
	return host, port, name, nil
}