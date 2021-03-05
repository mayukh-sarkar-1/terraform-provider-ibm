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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ----------------  Retrieve a reserved IP ---------------------------- //

// Define all the constants that matches with the given terrafrom attribute
const (
	// Request Param Constants
	isSubNetID     = "subnet_id" // This is subnet id for the given reserved ip
	isReservedIPID = "id"        // This is the id for the given reserved ip

	// Response Param Constants
	isReservedIP           = "address"     // This is the IP address for the given `isReservedIPID`
	isReservedIPAutoDelete = "auto_delete" // This mentions if reserved ip shall be deleted automatically
	isReservedIPCreatedAt  = "created_at"  // Time when the reserved ip was created
	isReservedIPhref       = "href"        // The url for the reserved ip
	// // TODO: Check if I can reuse `isReservedIPID` somehow
	// isReservedIPIDret = "id"            // This is the same reserved ip ID returned that was passed (may be redundant)
	isReservedIPName  = "name"          // This is system generated name given to reserved ip ( can be funny sometimes..;-) )
	isReservedIPOwner = "owner"         // Owner of this reserve ip..may be the username who assigned it
	isReservedIPType  = "resource_type" // This is the resource type
)

/*
	DataSourceIBMISReservedIP is the function
	that is called when use data source in terraform
	to get the attributes of reserve IP. This fucntion
	is expoerted outside of this module. All these are
	part of the main ibm package but just for sake of
	uniformity, we capitalize the first letter of the
	function.
*/
func dataSourceIBMISReservedIP() *schema.Resource {
	/*
		Go VPC call to get the reserved IP

		options := vpcService.NewGetSubnetReservedIPOptions(subnet_id, id)
		rip, response, err := vpcService.GetSubnetReservedIP(options)
	*/
	return &schema.Resource{
		Read: dataSdataSourceIBMISReservedIPRead,
		Schema: map[string]*schema.Schema{
			/*
				Request Parameters
				==================
				These are mandatory req parameters
				DOC: https://test.cloud.ibm.com/apidocs/vpc#get-subnet-reserved-ip
			*/
			isSubNetID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The subnet identifier.",
			},
			isReservedIPID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The reserved IP identifier.",
			},

			/*
				Response Parameters
				===================
				All of these are computed and an user doesn't need to provide
				these from outside.

				DOC: https://test.cloud.ibm.com/apidocs/vpc?code=go#get-subnet-reserved-ip
			*/

			isReservedIP: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IP address",
			},
			isReservedIPAutoDelete: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If set to true, this reserved IP will be automatically deleted",
			},
			isReservedIPCreatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the reserved IP was created.",
			},
			isReservedIPhref: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this reserved IP.",
			},
			isReservedIPName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user-defined or system-provided name for this reserved IP.",
			},
			isReservedIPOwner: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The owner of a reserved IP, defining whether it is managed by the user or the provider.",
			},
			isReservedIPType: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
		},
	}
}

// dataSdataSourceIBMISReservedIPRead is used when the
// reserved IPs are read from the vpc
func dataSdataSourceIBMISReservedIPRead(d *schema.ResourceData, meta interface{}) error {

	// First get the session for the VPC go SDK
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	// Get the inputs mentioned in the terraform script (Request Params)
	subnetID := d.Get(isSubNetID).(string)
	reservedIPID := d.Get(isReservedIPID).(string)
	msg := fmt.Sprintf("Trying to retrive the reserved IP for Subnet id: %s and reserved ip ID: %s", subnetID, reservedIPID)
	fmt.Println(" ✅ " + "\033[35m" + msg + "\033[0m")

	// Create reserved ip options (Basically a combination of subnet id and reserved ip id)
	options := sess.NewGetSubnetReservedIPOptions(subnetID, reservedIPID)
	reserveIP, response, err := sess.GetSubnetReservedIP(options)

	// Now check errors and some iternal issues
	if err != nil || response == nil {
		return fmt.Errorf("Error fetching the reserved IP %s\n%s", err, response)
	}

	// Now set the terraform variables that we have got as response
	d.SetId(*reserveIP.ID)
	d.Set(isReservedIPAutoDelete, *reserveIP.AutoDelete)
	d.Set(isReservedIPCreatedAt, (*reserveIP.CreatedAt).String())
	d.Set(isReservedIPhref, *reserveIP.Href)
	d.Set(isReservedIPName, *reserveIP.Name)
	d.Set(isReservedIPOwner, *reserveIP.Owner)
	d.Set(isReservedIPType, *reserveIP.ResourceType)
	return nil // By default there should be no error
}
