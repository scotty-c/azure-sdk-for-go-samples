// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package authorization

import (
	"context"
	"flag"
	"log"
	"os"
	"testing"

	"github.com/Azure-Samples/azure-sdk-for-go-samples/graphrbac"
	"github.com/Azure-Samples/azure-sdk-for-go-samples/internal/config"
	"github.com/Azure-Samples/azure-sdk-for-go-samples/internal/util"
	"github.com/Azure-Samples/azure-sdk-for-go-samples/resources"
)

// TestMain sets up the environment and initiates tests.
func TestMain(m *testing.M) {
	var err error
	err = config.ParseEnvironment()
	if err != nil {
		log.Fatalf("failed to parse env: %v\n", err)
	}
	err = config.AddFlags()
	if err != nil {
		log.Fatalf("failed to parse env: %v\n", err)
	}
	flag.Parse()

	code := m.Run()
	os.Exit(code)
}

func ExampleAssignRole() {
	var groupName = config.GenerateGroupName("Authorization")
	config.SetGroupName(groupName)

	ctx := context.Background()
	defer resources.Cleanup(ctx)

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		util.PrintAndLog(err.Error())
	}

	list, err := ListRoleDefinitions(ctx, "roleName eq 'Contributor'")
	if err != nil {
		util.PrintAndLog(err.Error())
	}
	util.PrintAndLog("got role definitions list")

	var userID string
	user, err := graphrbac.GetCurrentUser(ctx)
	if err != nil {
		log.Printf("could not get object for current user: %v\n", err)
		log.Printf("using service principal ID instead")
		userID = config.ClientID()
	} else {
		userID = *user.ObjectID
	}

	groupRole, err := AssignRole(ctx, userID, *list.Values()[0].ID)
	if err != nil {
		util.PrintAndLog(err.Error())
	}
	util.PrintAndLog("role assigned with resource group scope")

	subscriptionRole, err := AssignRoleWithSubscriptionScope(
		ctx, userID, *list.Values()[0].ID)
	if err != nil {
		util.PrintAndLog(err.Error())
	}
	util.PrintAndLog("role assigned with subscription scope")

	if !config.KeepResources() {
		DeleteRoleAssignment(ctx, *groupRole.ID)
		if err != nil {
			util.PrintAndLog(err.Error())
		}

		DeleteRoleAssignment(ctx, *subscriptionRole.ID)
		if err != nil {
			util.PrintAndLog(err.Error())
		}
	}

	// Output:
	// got role definitions list
	// role assigned with resource group scope
	// role assigned with subscription scope
}
