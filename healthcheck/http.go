package healthcheck

import (
	"time"
	"net/http"
	"fmt"
)

type HealthCheck struct {
	Target string
	Healthy bool
	Period time.Duration

	onChange OnChangeFunc

	run bool
}

type OnChangeFunc func()

func NewHealthCheck(target string, period time.Duration) *HealthCheck {
	return &HealthCheck{
		Target: target,
		Period: period,
		Healthy: false,
		onChange: nil,
		run: true,
	}
}

func (c *HealthCheck) OnChange(fu OnChangeFunc) {
	c.onChange = fu
}

func (c *HealthCheck) Start() {

	client := http.Client{Timeout: 3 * time.Second}
	go func() {
		for {
			_, err := client.Get(c.Target)
			if (err != nil && c.Healthy){
				c.Healthy = false
				c.onChange()
			} else if (err == nil && !c.Healthy) {
				c.Healthy = true
				c.onChange()
			}

			fmt.Println("Performed health check: ", err == nil, c.Healthy)
			time.Sleep(c.Period)
		}
	}()
}

func (c *HealthCheck) Stop() {
	c.run = false
}
