package windns

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider allows making changes to Windows DNS server
// Utilises Powershell to connect to domain controller
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("USERNAME", nil),
				Description: "Username to connect to AD.",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("PASSWORD", nil),
				Description: "The password to connect to AD.",
			},
			"server": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SERVER", nil),
				Description: "The AD server to connect to.",
			},
			"usessl": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("USESSL", false),
				Description: "Whether or not to use HTTPS to connect to WinRM",
			},
			"usessh": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("USESSH", false),
				Description: "Whether or not to use SSH to connect to WinRM",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"windns": resourceWinDNSRecord(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	username := d.Get("username").(string)
	if username == "" {
		return nil, fmt.Errorf("The 'username' property was not specified.")
	}

	usessh := d.Get("usessh").(string)

	password := d.Get("password").(string)
	if password == "" && usessh == "0" {
		return nil, fmt.Errorf("The 'password' property was not specified and usessh was false.")
	}

	server := d.Get("server").(string)
	if server == "" {
		return nil, fmt.Errorf("The 'server' property was not specified.")
	}

	usessl := d.Get("usessl").(string)

	f, err := ioutil.TempFile("", "terraform-windns")
	lockfile := f.Name()
	err = f.Close()
	err = os.Remove(f.Name())

	client := DNSClient{
		username: username,
		password: password,
		server:   server,
		usessl:   usessl,
		usessh:   usessh,
		lockfile: lockfile,
	}

	return &client, err
}

type DNSClient struct {
	username string
	password string
	server   string
	usessl   string
	usessh   string
	lockfile string
}
