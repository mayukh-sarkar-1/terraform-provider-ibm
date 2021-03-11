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

func TestAccIBMISReservedIP_basic(t *testing.T) {
	// resName := "manual-reserve-i"
	vpcName := fmt.Sprintf("tfresip-vpc-%d", acctest.RandIntRange(10, 100))
	subnetName := fmt.Sprintf("tfresip-subnet-%d", acctest.RandIntRange(10, 100))
	resIPName := fmt.Sprintf("tfresip-reservedip-%d", acctest.RandIntRange(10, 100))
	terraformTag := "data.ibm_is_subnet_reserved_ip.data_resip1"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIBMISReservedIPdataSoruceConfig(vpcName, subnetName, resIPName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(terraformTag, isReservedIPName, resIPName),
				),
			},
		},
	})
}

func testAccIBMISReservedIPdataSoruceConfig(vpcName, subnetName, reservedIPName string) string {
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

		data "ibm_is_subnet_reserved_ip" "data_resip1" {
			subnet = ibm_is_subnet.subnet1.id
			reserved_ip = ibm_is_subnet_reserved_ip.resip1.reserved_ip
		}
      `, vpcName, subnetName, reservedIPName)
}
