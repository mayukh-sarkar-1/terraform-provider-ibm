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
	"time"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isReservedIPProvisioning     = "provisioning"
	isReservedIPProvisioningDone = "done"
	isReservedIP                 = "reserved_ip"
)

func resourceIBMISReservedIP() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMISReservedIPCreate,
		Read:     resourceIBMISReservedIPRead,
		Update:   resourceIBMISReservedIPUpdate,
		Delete:   resourceIBMISReservedIPDelete,
		Exists:   resourceIBMISReservedIPExists,
		Importer: &schema.ResourceImporter{},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		/*
			d.Set(isReservedIPAddress, *rip.Address)
			d.Set(isReservedIP, *rip.ID)
			d.Set(isSubNetID, subnetID)
			d.Set(isReservedIPAutoDelete, *rip.AutoDelete)
			d.Set(isReservedIPCreatedAt, *rip.CreatedAt)
			d.Set(isReservedIPhref, *rip.Href)
			d.Set(isReservedIPName, *rip.Name)
			d.Set(isReservedIPOwner, *rip.Owner)
			d.Set(isReservedIPType, *rip.ResourceType)
		*/
		Schema: map[string]*schema.Schema{
			/*
				Request Parameters
				==================
				These are mandatory req parameters
				DOC: https://test.cloud.ibm.com/apidocs/vpc#create-subnet-reserved-ip
			*/
			isSubNetID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The subnet identifier.",
			},
			isReservedIPAutoDelete: {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "If set to true, this reserved IP will be automatically deleted",
			},
			isReservedIPName: {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The user-defined or system-provided name for this reserved IP.",
			},
			/*
				Response Parameters
				===================
				All of these are computed and an user doesn't need to provide
				these from outside.

				DOC: https://test.cloud.ibm.com/apidocs/vpc#create-subnet-reserved-ip
			*/

			isReservedIPAddress: {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The user-defined or system-provided name for this reserved IP.",
			},
			isReservedIP: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the reserved IP.",
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

// resourceIBMISReservedIPCreate Creates a reserved IP given a subnet ID
func resourceIBMISReservedIPCreate(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	// Getting the subnet ID
	subnetID := d.Get(isSubNetID).(string)

	options := sess.NewCreateSubnetReservedIPOptions(subnetID)

	// Getting the name of the reserved IP if given
	nameStr := ""
	if name, ok := d.GetOk(isReservedIPName); ok {
		nameStr = name.(string)
	}
	if nameStr != "" {
		options.Name = &nameStr
	}

	// Setting the auto delete
	autoDeleteBool := false
	if autoDelete, ok := d.GetOk(isReservedIPAutoDelete); ok {
		autoDeleteBool = autoDelete.(bool)
	}
	options.AutoDelete = &autoDeleteBool

	// Now create the reserved IP
	rip, response, err := sess.CreateSubnetReservedIP(options)
	if err != nil || response == nil {
		return fmt.Errorf("Error creating the reserved IP: %s\n%s", err, response)
	} else {
		msg := fmt.Sprintf("Created reserved IP for Subnet id: %s", subnetID)
		fmt.Println(" ✅ " + "\033[35m" + msg + "\033[0m")
	}

	// Set id for the reserved IP
	reservedIPID := *rip.ID
	d.SetId(fmt.Sprintf("%s/%s", subnetID, reservedIPID))

	// Finally call read method to read the resouce, set the variables and return
	return resourceIBMISReservedIPRead(d, meta)
}

func resourceIBMISReservedIPRead(d *schema.ResourceData, meta interface{}) error {
	rip, err := get(d, meta)
	if err != nil {
		return err
	}

	allIDs, err := idParts(d.Id())
	subnetID := allIDs[0]

	if rip != nil {
		d.Set(isReservedIPAddress, *rip.Address)
		d.Set(isReservedIP, *rip.ID)
		d.Set(isSubNetID, subnetID)
		d.Set(isReservedIPAutoDelete, *rip.AutoDelete)
		d.Set(isReservedIPCreatedAt, (*rip.CreatedAt).String())
		d.Set(isReservedIPhref, *rip.Href)
		d.Set(isReservedIPName, *rip.Name)
		d.Set(isReservedIPOwner, *rip.Owner)
		d.Set(isReservedIPType, *rip.ResourceType)
	}
	return nil
}

func resourceIBMISReservedIPUpdate(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	name := "" // Primarily name is something we are patching

	if d.HasChange(isReservedIPName) {
		name = d.Get(isReservedIPName).(string)
	}
	allIDs, err := idParts(d.Id())
	subnetID := allIDs[0]
	reservedIPID := allIDs[1]

	if name != "" {
		options := &vpcv1.UpdateSubnetReservedIPOptions{
			SubnetID: &subnetID,
			ID:       &reservedIPID,
		}

		patch := &vpcv1.ReservedIPPatch{
			Name: &name,
		}

		reservedIPPatch, err := patch.AsPatch()
		if err != nil {
			return fmt.Errorf("Error updating the reserved IP %s", err)
		}

		options.ReservedIPPatch = reservedIPPatch

		_, response, err := sess.UpdateSubnetReservedIP(options)
		if err != nil {
			return fmt.Errorf("Error updating the reserved IP %s\n%s", err, response)
		}
	}
	return resourceIBMISReservedIPRead(d, meta)
}

func resourceIBMISReservedIPDelete(d *schema.ResourceData, meta interface{}) error {

	// First get the reserved IP and check if it exists or not
	rip, err := get(d, meta)
	if err != nil {
		return err
	}
	if err == nil && rip == nil {
		// If there is no such reserved IP, it can not be deleted
		return fmt.Errorf("Can not find a reserved IP")
	}
	// Now delete the reserved IP if found
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	allIDs, err := idParts(d.Id())
	subnetID := allIDs[0]
	reservedIPID := allIDs[1]
	deleteOptions := sess.NewDeleteSubnetReservedIPOptions(subnetID, reservedIPID)
	response, err := sess.DeleteSubnetReservedIP(deleteOptions)
	if err != nil {
		return err
	}
	if err == nil && response == nil {
		return fmt.Errorf("Error deleting the reserverd ip %s in subnet %s", reservedIPID, subnetID)
	}
	d.SetId("")
	return nil
}

func resourceIBMISReservedIPExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	rip, err := get(d, meta)
	if err != nil {
		return false, err
	}
	if err == nil && rip == nil {
		return false, nil
	}
	return true, nil
}

// get is a generic function that gets the reserved ip given subnet id and reserved ip
func get(d *schema.ResourceData, meta interface{}) (*vpcv1.ReservedIP, error) {
	sess, err := vpcClient(meta)
	if err != nil {
		return nil, err
	}
	allIDs, err := idParts(d.Id())
	subnetID := allIDs[0]
	reservedIPID := allIDs[1]
	options := sess.NewGetSubnetReservedIPOptions(subnetID, reservedIPID)
	rip, response, err := sess.GetSubnetReservedIP(options)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil, nil
		}
		return nil, fmt.Errorf("Error Getting Reserved IP : %s\n%s", err, response)
	}
	return rip, nil
}
