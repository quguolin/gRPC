package consul

import (
	"fmt"
	"net/http"
	"time"

	consulapi "github.com/hashicorp/consul/api"
)

type Config struct {
	IP   string
	Port int
	Tag  []string
	Name string
	//CheckPort health check port
	CheckUrl string
	CheckPort int
	CheckInter int
	CheckDeReg int
}

type Consul struct {

}

func (con *Consul)Register(c *Config) (error) {
	if c == nil{
		return fmt.Errorf("config can not null")
	}
	if c.Name == "" || c.IP == "" || c.Port == 0{
		return fmt.Errorf("name or ip or port can not null")
	}
	config := consulapi.DefaultConfig()
	client, err := consulapi.NewClient(config)
	if err != nil {
		return err
	}
	reg := &consulapi.AgentServiceRegistration{
		ID:      fmt.Sprintf("%s-%s-%d", c.Name, c.IP, c.Port),
		Address:c.IP,
		Port:c.Port,
		Name:c.Name,
		Tags:c.Tag,
	}
	if c.CheckUrl != "" && c.CheckPort != 0{
		inter := time.Duration(c.CheckInter) * time.Second
		deReg := time.Duration(c.CheckDeReg) * time.Second
		reg.Check = &consulapi.AgentServiceCheck{ 				// 健康检查
			Interval:                       inter.String(),		// 健康检查间隔
			HTTP:                           fmt.Sprintf("http://%s:%d%s", c.IP, c.CheckPort, c.CheckUrl), // 执行健康检查的地址
			DeregisterCriticalServiceAfter: deReg.String(),    	// check失败后 删除本服务
		}
	}
	err = client.Agent().ServiceRegister(reg)
	if err != nil {
		return err
	}
	http.HandleFunc(c.CheckUrl, func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintln(writer, "consulCheck")
	})
	if err = http.ListenAndServe(fmt.Sprintf(":%d", c.CheckPort), nil); err != nil {
		return err
	}
	return nil
}