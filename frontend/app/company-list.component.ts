import { Component, Input, SimpleChange,  ViewChild } from '@angular/core';

import { BaseChartDirective } from 'ng2-charts/ng2-charts';

import { Developer } from './models'

@Component({
  selector: 'company-list',
  template: `
  <div class="container-fluid">
    <div class="row">
    <div class="col-md-6 col-xs-12">
        <h2>Company statistics</h2>
        <table class="table table-striped table-bordered">
            <thead>
            <tr>
                <th>Name</th>
                <th>PR Count</th>
            </tr>
            </thead>
            <tbody>
            <tr *ngFor="let i of comps">
                <td>{{i.name}}</td>
                <td>{{i.pr_count}}</td>
            </tr>
            </tbody>
        </table>
    </div>
    <div class="col-md-6 col-xs-12">
        <div style="display: block">
          <canvas baseChart
                  [data]="pieChartData"
                  [labels]="pieChartLabels"
                  [chartType]="pieChartType">
          </canvas>
        </div>
    </div>
    </div>
  </div>
  `
})

export class CompanyListComponent{
  @Input()
  comps: Developer[];

  pieChartData = [1,2,3];
  pieChartLabels: Array<string> = [];
  pieChartType = "doughnut";

  @ViewChild(BaseChartDirective) chart: BaseChartDirective;

  ngOnChanges(changes: {[propKey: string]: SimpleChange}) {
    let log: string[] = [];
    for (let propName in changes) {
      if(propName == "comps") {
          let changedProp = changes[propName];
          let to = changedProp.currentValue;

          if(to != null) {
            this.pieChartLabels = to.map(comp => comp.name);
            this.pieChartData = to.map(comp => comp.pr_count);
            //not angulary, hacky, but works
            this.chart.labels = this.pieChartLabels;
            this.chart.data = this.pieChartData;
            this.chart.ngOnChanges({});
          }
      }
    }
  }
}

