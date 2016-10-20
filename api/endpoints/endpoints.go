// Copyright 2016 Mirantis
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package endpoints

import (
	"github.com/gin-gonic/gin"

	"github.com/Mirantis/statkube/db"
	"github.com/Mirantis/statkube/models"
)

func GetPRStatsDev(c *gin.Context) {
	db := db.GetDB()
	defer db.Close()

	devs, err := models.GetDevStats(db)
	if err != nil {
		c.Error(err)
		c.JSON(500, []int{})
	}

	c.JSON(200, devs)
}

func GetPRStatsCompany(c *gin.Context) {
	db := db.GetDB()
	defer db.Close()

	devs, err := models.GetCompanyStats(db)
	if err != nil {
		c.Error(err)
		c.JSON(500, []int{})
	}

	c.JSON(200, devs)
}
