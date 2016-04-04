// Copyright 2015 ThoughtWorks, Inc.

// This file is part of Gauge.

// Gauge is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// Gauge is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with Gauge.  If not, see <http://www.gnu.org/licenses/>.

package result

import (
	"github.com/getgauge/gauge/gauge_messages"
	"github.com/golang/protobuf/proto"
)

type SpecResult struct {
	ProtoSpec            *gauge_messages.ProtoSpec
	ScenarioFailedCount  int
	ScenarioCount        int
	IsFailed             bool
	FailedDataTableRows  []int32
	ExecutionTime        int64
	Skipped              bool
	ScenarioSkippedCount int
}

func (specResult *SpecResult) SetFailure() {
	specResult.IsFailed = true
}

func (specResult *SpecResult) AddSpecItems(resolvedItems []*gauge_messages.ProtoItem) {
	specResult.ProtoSpec.Items = append(specResult.ProtoSpec.Items, resolvedItems...)
}

func (specResult *SpecResult) AddScenarioResults(scenarioResults []*ScenarioResult) {
	for _, scenarioResult := range scenarioResults {
		if scenarioResult.ProtoScenario.GetFailed() {
			specResult.IsFailed = true
			specResult.ScenarioFailedCount++
		}
		specResult.AddExecTime(scenarioResult.ProtoScenario.GetExecutionTime())
		specResult.ProtoSpec.Items = append(specResult.ProtoSpec.Items, &gauge_messages.ProtoItem{ItemType: gauge_messages.ProtoItem_Scenario.Enum(), Scenario: scenarioResult.ProtoScenario})
	}
	specResult.ScenarioCount += len(scenarioResults)
}

func (specResult *SpecResult) AddTableDrivenScenarioResult(scenarioResults [][](*ScenarioResult)) {
	numberOfScenarios := len(scenarioResults[0])

	for scenarioIndex := 0; scenarioIndex < numberOfScenarios; scenarioIndex++ {
		protoTableDrivenScenario := &gauge_messages.ProtoTableDrivenScenario{Scenarios: make([]*gauge_messages.ProtoScenario, 0)}
		scenarioFailed := false
		for rowIndex, eachRow := range scenarioResults {
			protoScenario := eachRow[scenarioIndex].ProtoScenario
			protoTableDrivenScenario.Scenarios = append(protoTableDrivenScenario.GetScenarios(), protoScenario)
			specResult.AddExecTime(protoScenario.GetExecutionTime())
			if protoScenario.GetFailed() {
				scenarioFailed = true
				specResult.FailedDataTableRows = append(specResult.FailedDataTableRows, int32(rowIndex))
			}
		}
		if scenarioFailed {
			specResult.ScenarioFailedCount++
			specResult.IsFailed = true
		}
		protoItem := &gauge_messages.ProtoItem{ItemType: gauge_messages.ProtoItem_TableDrivenScenario.Enum(), TableDrivenScenario: protoTableDrivenScenario}
		specResult.ProtoSpec.Items = append(specResult.ProtoSpec.Items, protoItem)
	}
	specResult.ProtoSpec.IsTableDriven = proto.Bool(true)
	specResult.ScenarioCount += numberOfScenarios
}

func (specResult *SpecResult) AddExecTime(execTime int64) {
	specResult.ExecutionTime += execTime
}

func (specResult *SpecResult) getPreHook() **(gauge_messages.ProtoHookFailure) {
	return &specResult.ProtoSpec.PreHookFailure
}

func (specResult *SpecResult) getPostHook() **(gauge_messages.ProtoHookFailure) {
	return &specResult.ProtoSpec.PostHookFailure
}

func (specResult *SpecResult) setFileName(fileName string) {
	specResult.ProtoSpec.FileName = proto.String(fileName)
}
