// Copyright (c) 2022, Salesforce, Inc.
// All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// For full license text, see the LICENSE file in the repo root or https://opensource.org/licenses/BSD-3-Clause

package service

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mce.salesforce.com/sprinkler/database"
	"mce.salesforce.com/sprinkler/database/table"
	"mce.salesforce.com/sprinkler/model"
)

type Control struct {
	db *gorm.DB
}

type postWorkflowReq struct {
	Name        string    `json:"name" binding:"required"`
	Artifact    string    `json:"artifact" binding:"required"`
	Command     string    `json:"command" binding:"required"`
	Every       string    `json:"every" binding:"required"`
	NextRuntime time.Time `json:"nextRuntime" binding:"required"`
	Backfill    bool      `json:"backfill" binding:"required"`
	Owner       *string   `json:"owner"`
	IsActive    bool      `json:"isActive" binding:"required"`
}

func NewControl() *Control {
	return &Control{
		db: database.GetInstance(),
	}
}

func (ctrl *Control) postWorkflow(c *gin.Context) {
	var body postWorkflowReq
	if err := c.BindJSON(&body); err != nil {
		// bad request
		fmt.Println(err)
		return
	}

	every, err := model.ParseEvery(body.Every)
	if err != nil {
		fmt.Println(err)
		return
	}

	wf := table.Workflow{
		Name:        body.Name,
		Artifact:    body.Artifact,
		Command:     body.Command,
		Every:       every,
		NextRuntime: body.NextRuntime,
		Backfill:    body.Backfill,
		Owner:       body.Owner,
		IsActive:    body.IsActive,
	}
	ctrl.db.Create(&wf)
	c.JSON(http.StatusOK, "OK")
}

func (ctrl *Control) Run() {
	r := gin.Default()
	r.POST("v1/workflow", ctrl.postWorkflow)
	r.Run()
}
