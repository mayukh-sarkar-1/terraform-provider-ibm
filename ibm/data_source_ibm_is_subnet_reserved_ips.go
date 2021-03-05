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

// ----------------  Retrieve all the reserved IP ---------------------------- //

// Define all the constants that matches with the given terrafrom attribute
const (
	// Request Param Constants
	isReservedIPLimit  = "limit"        // Number of reserved ips to return
	isReservedIPSort   = "sort"         // Attribute on which the IPs are sorted
	isReservedIPs      = "reserved_ips" // List of all the reserved IPs
	isReservedIPsCount = "total_count"  // Total number of reserved IPs
)

func dataSourceIBMISReservedIPs() *schema.Resource {
	return &schema.Resource{
		Read: dataSdataSourceIBMISReservedIPsRead,
		Schema: map[string]*schema.Schema{
			/*
				Request Parameters
				==================
				These are mandatory req parameters
				DOC: https://test.cloud.ibm.com/apidocs/vpc#list-subnet-reserved-ips
			*/
			isSubNetID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The subnet identifier.",
			},
			isReservedIPLimit: {
				Type:        schema.TypeInt,
				Default:     50,
				Optional:    true,
				Description: "The number of resources to return on a page",
			},
			isReservedIPSort: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "address",
				Description: "Sorts the returned collection by the specified field name in ascending order",
			},
			/*
				Response Parameters
				===================
				All of these are computed and an user doesn't need to provide
				these from outside.

				DOC: https://test.cloud.ibm.com/apidocs/vpc#list-subnet-reserved-ips
			*/

			isReservedIPs: {
				Type:        schema.TypeList,
				Description: "Collection of reserved IPs in this subnet.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isReservedIP: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP address",
						},
						isReservedIPAutoDelete: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "If reserved ip shall be deleted automatically",
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
						isReservedIPID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this reserved IP",
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
				},
			},
			isReservedIPsCount: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of resources across all pages",
			},
		},
	}
}

func dataSdataSourceIBMISReservedIPsRead(d *schema.ResourceData, meta interface{}) error {
	// First get the session for the VPC go SDK
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	// Get the inputs mentioned in the terraform script (Request Params)
	subnetID := d.Get(isSubNetID).(string)
	limit := d.Get(isReservedIPLimit).(int)
	sortKey := d.Get(isReservedIPSort).(string)
	msg := fmt.Sprintf("Trying to retrive all the reserved IPs for Subnet id: %s", subnetID)
	fmt.Println(" ✅ " + "\033[35m" + msg + "\033[0m")

	// Create reserved ip options (Passing only the subnet_id)
	options := sess.NewListSubnetReservedIpsOptions(subnetID)

	// Now get all the reserved IPs
	Ips, Limit, err := getAllReservedIPs(sess, options, limit, sortKey, subnetID)
	if err != nil {
		return fmt.Errorf("❌ Error getting reserved IPs for subnet %s", subnetID)
	}
	// Now set the outputs in terraform
	d.Set(isReservedIPs, *Ips)
	d.Set(isReservedIPLimit, Limit)
	d.Set(isReservedIPsCount, len(*Ips))
	d.SetId(time.Now().UTC().String()) // This is not any reserved ip or subnet id but state id
	return nil
}

func getAllReservedIPs(sess *vpcv1.VpcV1, options *vpcv1.ListSubnetReservedIpsOptions,
	limit int, subnetID, sortKey string) (*[]map[string]interface{}, int, error) {
	// Add other options
	options.SetLimit(int64(limit))
	options.SetSort(sortKey)

	allResIPs := []vpcv1.ReservedIP{}
	l := -1
	start := ""
	for {
		if start != "" {
			options.Start = &start
		}
		// Get the reserved IP iterator
		reserveIPs, response, err := sess.ListSubnetReservedIps(options)
		// Now check errors and some iternal issues
		if err != nil || response == nil {
			return nil, l, fmt.Errorf("Error fetching the reserved IP %s\n%s", err, response)
		}

		// Get the first page from the iterator
		start = GetNext(reserveIPs.Next)
		l = int(*reserveIPs.Limit)
		allResIPs = append(allResIPs, reserveIPs.ReservedIps...)
		if start == "" {
			break
		}
	}

	// Now load all the reserved IPs into a map
	reservedIPs := make([]map[string]interface{}, 0)
	for _, data := range allResIPs {
		// ips := data.(*vpcv1.ReservedIP)
		reservedIP := map[string]interface{}{
			isReservedIPID:         *data.ID,
			isReservedIPName:       *data.Name,
			isReservedIP:           *data.Address,
			isReservedIPAutoDelete: *data.AutoDelete,
			isReservedIPCreatedAt:  (*data.CreatedAt).String(),
			isReservedIPhref:       *data.Href,
			isReservedIPOwner:      *data.Owner,
			isReservedIPType:       *data.ResourceType,
		}
		reservedIPs = append(reservedIPs, reservedIP)
	}
	return &reservedIPs, l, nil
}
