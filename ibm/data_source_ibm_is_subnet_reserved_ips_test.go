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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIBMISReservedIPs_basic(t *testing.T) {
	resName := "data.ibm_is_subnet_reserved_ips.test_res_ip"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIBMISReservedIPSdataSoruceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resName, "reserved_ips.0.address"),
				),
			},
		},
	})
}

func testAccIBMISReservedIPSdataSoruceConfig() string {
	// status filter defaults to empty
	return fmt.Sprintf(`

      data "ibm_is_subnet_reserved_ips" "test_res_ip" {
      	subnet_id = "0716-d335ad68-1538-4d9f-8bc4-04c745f662c2"
      }
      `)
}
