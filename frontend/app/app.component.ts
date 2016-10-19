import { Component } from '@angular/core';
import * as moment from 'moment/moment';

@Component({
  selector: 'my-app',
  template: `
    <h1>STATKUBE</h1>
    <div class="form-group">
        <select [ngModel]="selected_period" class="form-control" (ngModelChange)="onPeriodChange($event)">
            <option [ngValue]="time" *ngFor="let time of times">{{time.name}}</option>
        </select>
    </div>
    <span>start: {{selected_time.start}} end: {{selected_time.end}}</span>
    <div class="form-group">
        <label>from</label>
        <input [(ngModel)]="selected_time.start" type="date" class="form-control"/>
    </div>
    <div class="form-group">
        <label>to</label>
        <input [(ngModel)]="selected_time.end" type="date" class="form-control"/>
    </div>
    <h2>Engineer statistics</h2>
    <table class="table table-striped table-bordered">
        <thead>
        <tr>
            <th>Name</th>
            <th>PR Count</th>
        </tr>
        </thead>
        <tbody>
        <tr *ngFor="let dev of devs">
            <td>{{dev.name}}</td>
            <td>{{dev.pr_count}}</td>
        </tr>
        </tbody>
    </table>
  `
})

export class AppComponent {
  devs = DEVS;
  selected_time = new TimeConstraint("", "", "");
  selected_period = new TimeConstraint("", "", "");
  times = TIMES_SELECTABLE;

  onPeriodChange(event: TimeConstraint): void {
	  this.selected_time.start = event.start;
	  this.selected_time.end = event.end;
	  this.selected_time.name = event.name;
  };
}

export class Developer {
  name: string;
  pr_count: number;
  constructor(name: string, pr_count: number) {
    this.name = name;
    this.pr_count = pr_count;
  }
}

export class TimeConstraint {
    start: string;
    end: string;
    name: string;
    constructor(start: string, end: string, name: string) {
        this.start = start;
        this.end = end;
        this.name = name;
    }
}

const DEVS: Developer[] = [
    new Developer("super dev", 12),
    new Developer("super dev2", 13),
    new Developer("super dev3", 5),
    new Developer("super dev4", 120),
];

const MOMENT_FORMAT = "YYYY-MM-DD";

const TIMES_SELECTABLE: TimeConstraint[] = [
    new TimeConstraint("", "", ""),
    new TimeConstraint(moment().startOf("week").format(MOMENT_FORMAT), moment().format(MOMENT_FORMAT), "this week"),
    new TimeConstraint(moment().startOf("week").subtract(1, 'week').format(MOMENT_FORMAT), moment().startOf("week").format(MOMENT_FORMAT), "last week"),
]
