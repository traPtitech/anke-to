package controller

import (
	"context"
	"slices"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/gofrs/uuid"
	"github.com/traPtitech/anke-to/model"
	"github.com/traPtitech/anke-to/traq"
)

func isAllTargetsReponded(targets []model.Targets, respondents []model.Respondents) bool {
	respondentsString := []string{}
	for _, respondent := range respondents {
		respondentsString = append(respondentsString, respondent.UserTraqid)
	}

	for _, target := range targets {
		if !slices.Contains(respondentsString, target.UserTraqid) {
			return false
		}
	}
	return true
}

func minimizeUsersAndGroups(users []string, groups []uuid.UUID) ([]string, []uuid.UUID, error) {
	ctx := context.Background()
	client := traq.NewTraqAPIClient()
	userSet := mapset.NewSet[string]()
	for _, user := range users {
		userSet.Add(user)
	}
	groupUserSet := mapset.NewSet[string]()
	for _, group := range groups {
		members, err := client.GetGroupMembers(ctx, group.String())
		if err != nil {
			return nil, nil, err
		}
		for _, member := range members {
			memberTraqID, err := client.GetUserTraqID(ctx, member.Id)
			if err != nil {
				return nil, nil, err
			}
			groupUserSet.Add(memberTraqID)
		}
	}
	userSet = userSet.Difference(groupUserSet)
	return userSet.ToSlice(), groups, nil
}

func rollOutUsersAndGroups(users []string, groups []string) ([]string, error) {
	ctx := context.Background()
	client := traq.NewTraqAPIClient()
	userSet := mapset.NewSet[string]()
	for _, user := range users {
		userSet.Add(user)
	}
	for _, group := range groups {
		members, err := client.GetGroupMembers(ctx, group)
		if err != nil {
			return nil, err
		}
		for _, member := range members {
			memberTraqID, err := client.GetUserTraqID(ctx, member.Id)
			if err != nil {
				return nil, err
			}
			userSet.Add(memberTraqID)
		}
	}
	return userSet.ToSlice(), nil
}

func uuid2GroupNames(users []string) ([]string, error) {
	ctx := context.Background()
	client := traq.NewTraqAPIClient()
	groupNames := []string{}
	for _, user := range users {
		groupName, err := client.GetGroupName(ctx, user)
		if err != nil {
			return nil, err
		}
		groupNames = append(groupNames, groupName)
	}
	return groupNames, nil
}
