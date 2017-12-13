package vcd

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	govcd "github.com/ukcloud/govcloudair"
)

func TestAccVcdVApp(t *testing.T) {
	var vapp govcd.VApp

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVcdVAppDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccResourceVAppBase(testAccResourceVMSimple()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVcdVAppExists("vcd_vapp.foobar", &vapp),
					resource.TestCheckResourceAttr(
						"vcd_vapp.foobar", "name", os.Getenv("VCD_VAPP_NAME")),
					// resource.TestCheckResourceAttr(
					// 	"vcd_vapp.foobar", "ip", "10.10.102.160"),
					// resource.TestCheckResourceAttr(
					// 	"vcd_vapp.foobar", "power_on", "true"),
				),
			},

			// resource.TestStep{
			// 	Config: fmt.Sprintf(testAccCheckVcdVApp_basic, os.Getenv("VCD_EDGE_GATEWAY"), os.Getenv("VCD_EDGE_GATEWAY")),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		resource.TestCheckResourceAttr(
			// 			"vcd_vapp.foobar_allocated", "name", "foobar-allocated"),
			// 		resource.TestCheckResourceAttr(
			// 			"vcd_vapp.foobar_allocated", "ip", "allocated"),
			// 		resource.TestCheckResourceAttr(
			// 			"vcd_vapp.foobar_allocated", "power_on", "true"),
			// 	),
			// },

			// resource.TestStep{
			// 	Config: fmt.Sprintf(testAccCheckVcdVApp_powerOff, os.Getenv("VCD_EDGE_GATEWAY")),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckVcdVAppExists("vcd_vapp.foobar", &vapp),
			// 		testAccCheckVcdVAppAttributes_off(&vapp),
			// 		resource.TestCheckResourceAttr(
			// 			"vcd_vapp.foobar", "name", "foobar"),
			// 		resource.TestCheckResourceAttr(
			// 			"vcd_vapp.foobar", "ip", "10.10.103.160"),
			// 		resource.TestCheckResourceAttr(
			// 			"vcd_vapp.foobar", "power_on", "false"),
			// 	),
			// },
		},
	})
}

func testAccCheckVcdVAppDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*VCDClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "foobar" {
			continue
		}

		_, err := conn.OrgVdc.FindVAppByName(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("VPCs still exist")
		}

		return nil
	}

	return nil
}

func testAccCheckVcdVAppExists(n string, vapp *govcd.VApp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VAPP ID is set")
		}

		conn := testAccProvider.Meta().(*VCDClient)

		resp, err := conn.OrgVdc.FindVAppByName(rs.Primary.ID)
		if err != nil {
			return err
		}

		*vapp = resp

		return nil
	}
}

func testAccResourceVAppBase(vms string) string {
	return fmt.Sprintf(`
resource "vcd_vapp" "test-vapp" {
  name     = "%s"
  power_on = "0"

  network = [
    "%s",
    "%s",
  ]

  %s
}
`,
		os.Getenv("VCD_VAPP_NAME"),
		os.Getenv("VCD_VAPP_NETWORK_1"),
		os.Getenv("VCD_VAPP_NETWORK_2"),
		vms,
	)
}

func testAccResourceVMSimple() string {
	return fmt.Sprintf(`vm = {
  name          = "%s"
  catalog_name  = "%s"
  template_name = "%s"
  memory        = 1024
  cpus          = 
  network = {
    name               = "%s"
    ip_allocation_mode = "POOL"
    is_primary         = true
    adapter_type       = "VMXNET3"
  
}`,
		"testSimple",
		os.Getenv("VCD_VAPP_NETWORK_1"),
		os.Getenv("VCD_CATALOG"),
		os.Getenv("VCD_TEMPLATE"))
}

func testAccResourceVMSimpleStep2() string {
	return fmt.Sprintf(`vm = {
  name          = "%s"
  catalog_name  = "%s"
  template_name = "%s"
  memory        = 2048
  cpus          = 
  network = {
    name               = "%s"
    ip_allocation_mode = "POOL"
    is_primary         = true
    adapter_type       = "VMXNET3"
  
}`,
		"testSimple",
		os.Getenv("VCD_VAPP_NETWORK_1"),
		os.Getenv("VCD_CATALOG"),
		os.Getenv("VCD_TEMPLATE"))
}

// const testAccCheckVcdVApp_basic = `
// resource "vcd_network" "foonet" {
// 	name = "foonet"
// 	edge_gateway = "%s"
// 	gateway = "10.10.102.1"
// 	static_ip_pool {
// 		start_address = "10.10.102.2"
// 		end_address = "10.10.102.254"
// 	}
// }

// resource "vcd_network" "foonet3" {
// 	name = "foonet3"
// 	edge_gateway = "%s"
// 	gateway = "10.10.202.1"
// 	static_ip_pool {
// 		start_address = "10.10.202.2"
// 		end_address = "10.10.202.254"
// 	}
// }

// resource "vcd_vapp" "foobar" {
//   name          = "foobar"
//   template_name = "Skyscape_CentOS_6_4_x64_50GB_Small_v1.0.1"
//   catalog_name  = "Skyscape Catalogue"
//   network_name  = "${vcd_network.foonet.name}"
//   memory        = 1024
//   cpus          = 1
//   ip            = "10.10.102.160"
// }

// resource "vcd_vapp" "foobar_allocated" {
//   name          = "foobar-allocated"
//   template_name = "Skyscape_CentOS_6_4_x64_50GB_Small_v1.0.1"
//   catalog_name  = "Skyscape Catalogue"
//   network_name  = "${vcd_network.foonet3.name}"
//   memory        = 1024
//   cpus          = 1
//   ip            = "allocated"
// }
// `

// const testAccCheckVcdVApp_powerOff = `
// resource "vcd_network" "foonet2" {
// 	name = "foonet2"
// 	edge_gateway = "%s"
// 	gateway = "10.10.103.1"
// 	static_ip_pool {
// 		start_address = "10.10.103.2"
// 		end_address = "10.10.103.170"
// 	}

// 	dhcp_pool {
// 		start_address = "10.10.103.171"
// 		end_address = "10.10.103.254"
// 	}
// }

// resource "vcd_vapp" "foobar" {
//   name          = "foobar"
//   template_name = "Skyscape_CentOS_6_4_x64_50GB_Small_v1.0.1"
//   catalog_name  = "Skyscape Catalogue"
//   network_name  = "${vcd_network.foonet2.name}"
//   memory        = 1024
//   cpus          = 1
//   ip            = "10.10.103.160"
//   power_on      = false
// }
// `

// func testAccCheckVcdVAppAttributes(vapp *govcd.VApp) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {

// 		if vapp.VApp.Name != "foobar" {
// 			return fmt.Errorf("Bad name: %s", vapp.VApp.Name)
// 		}

// 		if vapp.VApp.Name != vapp.VApp.Children.VM[0].Name {
// 			return fmt.Errorf("VApp and VM names do not match. %s != %s",
// 				vapp.VApp.Name, vapp.VApp.Children.VM[0].Name)
// 		}

// 		status, _ := vapp.GetStatus()
// 		if status != "POWERED_ON" {
// 			return fmt.Errorf("VApp is not powered on")
// 		}

// 		return nil
// 	}
// }

// func testAccCheckVcdVAppAttributes_off(vapp *govcd.VApp) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {

// 		if vapp.VApp.Name != "foobar" {
// 			return fmt.Errorf("Bad name: %s", vapp.VApp.Name)
// 		}

// 		if vapp.VApp.Name != vapp.VApp.Children.VM[0].Name {
// 			return fmt.Errorf("VApp and VM names do not match. %s != %s",
// 				vapp.VApp.Name, vapp.VApp.Children.VM[0].Name)
// 		}

// 		status, _ := vapp.GetStatus()
// 		if status != "POWERED_OFF" {
// 			return fmt.Errorf("VApp is still powered on")
// 		}

// 		return nil
// 	}
// }
