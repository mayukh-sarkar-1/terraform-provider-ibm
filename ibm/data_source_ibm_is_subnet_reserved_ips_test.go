/* IBM Confidential
*  Object Code Only Source Materials
*  5747-SM3
*  (c) Copyright IBM Corp. 2017,2021
*
*  The source code for this program is not published or otherwise divested
*  of its trade secrets, irrespective of what has been deposited with the
*  U.S. Copyright Office.
 */

package ibm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIBMISReservedIPs_basic(t *testing.T) {
	terraformTag := "data.ibm_is_subnet_reserved_ips.data_resip"
	vpcName := fmt.Sprintf("tfresip-vpc-%d", acctest.RandIntRange(10, 100))
	subnetName := fmt.Sprintf("tfresip-subnet-%d", acctest.RandIntRange(10, 100))
	reservedIPName := fmt.Sprintf("tfresip-reservedip-%d", acctest.RandIntRange(10, 100))
	reservedIPName2 := fmt.Sprintf("tfresip-reservedip-%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIBMISReservedIPSdataSoruceConfig(vpcName, subnetName, reservedIPName, reservedIPName2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(terraformTag, "reserved_ips.0.address"),
					// resource.TestCheckResourceAttrSet(terraformTag, "reserved_ips.1.address"),
				),
			},
		},
	})
}

func testAccIBMISReservedIPSdataSoruceConfig(vpcName, subnetName, reservedIPName, reservedIPName2 string) string {
	// status filter defaults to empty
	return fmt.Sprintf(`
		resource "ibm_is_vpc" "vpc1" {
			name = "%s"
		}

		resource "ibm_is_subnet" "subnet1" {
			name                     = "%s"
			vpc                      = ibm_is_vpc.vpc1.id
			zone                     = "us-south-1"
			total_ipv4_address_count = 256
		}

		resource "ibm_is_subnet_reserved_ip" "resip1" {
			subnet = ibm_is_subnet.subnet1.id
			name = "%s"
		}
		
		resource "ibm_is_subnet_reserved_ip" "resip2" {
			subnet = ibm_is_subnet.subnet1.id
			name = "%s"
		}
        
		data "ibm_is_subnet_reserved_ips" "data_resip" {
      	  subnet = ibm_is_subnet.subnet1.id
        }
      `, vpcName, subnetName, reservedIPName, reservedIPName2)
}
