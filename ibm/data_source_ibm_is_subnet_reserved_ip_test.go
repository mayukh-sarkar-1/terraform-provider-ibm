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

func TestAccIBMISReservedIP_basic(t *testing.T) {
	// resName := "manual-reserve-i"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIBMISReservedIPdataSoruceConfig(),
				// Check:  resource.ComposeTestCheckFunc(
				// // resource.TestCheckResourceAttr("data.ibm_is_subnet_reserved_ip.test_res_ip", "name", reservedIPName),
				// ),
			},
		},
	})
}

func testAccIBMISReservedIPdataSoruceConfig() string {
	// status filter defaults to empty
	return fmt.Sprintf(`
      data "ibm_is_subnet_reserved_ip" "test_res_ip" {
      	subnet_id = "0716-d335ad68-1538-4d9f-8bc4-04c745f662c2"
      	id = "0716-7de74d73-4686-41ec-8e65-c84a173a75bf"
      }
      `)
}
