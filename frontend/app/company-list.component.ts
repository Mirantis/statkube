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
        <ng-table [config]="config"
          (tableChanged)="onChangeTable(config)"
          [rows]="rows" [columns]="columns">
        </ng-table>
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
  public comps: Developer[];
  public rows: Developer[];

  public columns:Array<any> = [
      {title: 'Name', name: 'name', filtering: {filterString: '', placeholder: 'Filter by name'}},
      {title: 'PR count', name: 'pr_count', sort: 'desc'},
  ];

  public config:any = {
      paging: false,
      sorting: {columns: this.columns},
      filtering: {filterString: ''},
      className: ['table-striped', 'table-bordered']
  };

  public onChangeTable(config:any):any {
      if (config.filtering) {
          Object.assign(this.config.filtering, config.filtering);
      }

      if (config.sorting) {
          Object.assign(this.config.sorting, config.sorting);
      }

      let filteredData = this.changeFilter(this.comps, this.config);
      let sortedData = this.changeSort(filteredData, this.config);
      this.rows = sortedData
  };

  public changeSort(data:any, config:any):any {
      if (!config.sorting) {
          return data;

      }

      let columns = this.config.sorting.columns || [];
      let columnName:string = void 0;
      let sort:string = void 0;

      for (let i = 0; i < columns.length; i++) {
          if (columns[i].sort !== '' && columns[i].sort !== false) {
              columnName = columns[i].name;
              sort = columns[i].sort;

          }

      }

      if (!columnName) {
          return data;

      }

      return data.sort((previous:any, current:any) => {
          if (previous[columnName] > current[columnName]) {
              return sort === 'desc' ? -1 : 1;
          } else if (previous[columnName] < current[columnName]) {
              return sort === 'asc' ? -1 : 1;
          }
          return 0;
      });
  }

  public changeFilter(data:any, config:any):any {
          let filteredData:Array<any> = data;
          this.columns.forEach((column:any) => {
              if (column.filtering) {
                  filteredData = filteredData.filter((item:any) => {
                                return item[column.name].match(column.filtering.filterString);
                  });
              }
          });

          if (!config.filtering) {
                    return filteredData;
          }

          if (config.filtering.columnName) {
                    return filteredData.filter((item:any) =>
                                                      item[config.filtering.columnName].match(this.config.filtering.filterString));
          }

              let tempArray:Array<any> = [];
              filteredData.forEach((item:any) => {
                        let flag = false;
                        this.columns.forEach((column:any) => {
                            if (item[column.name].toString().match(this.config.filtering.filterString)) {
                                          flag = true;
                            }
                        });
                        if (flag) {
                                    tempArray.push(item);
                        }
              });
                  filteredData = tempArray;

                      return filteredData;
  }

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
            this.onChangeTable(this.config);

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

