import { Component, Input } from '@angular/core';
import * as moment from 'moment/moment';
import { TimeConstraint } from './models'

@Component({
	  selector: 'time-picker',
	  template: `
	  <div>
		<div class="form-group">
			<select [ngModel]="selected_period" class="form-control" (ngModelChange)="onPeriodChange($event)">
				<option [ngValue]="time" *ngFor="let time of times">{{time.name}}</option>
			</select>
		</div>
		<div class="form-group">
			<label>from</label>
			<input [(ngModel)]="selected_time.start" type="date" class="form-control"/>
		</div>
		<div class="form-group">
			<label>to</label>
			<input [(ngModel)]="selected_time.end" type="date" class="form-control"/>
		</div>
	  </div>
	  `
})

export class TimePickerComponent{
  @Input()
  selected_time: TimeConstraint;

  selected_period = new TimeConstraint("", "", "");
  times = TIMES_SELECTABLE;

  onPeriodChange(event: TimeConstraint): void {
	  this.selected_time.start = event.start;
	  this.selected_time.end = event.end;
	  this.selected_time.name = event.name;
  };
}

const MOMENT_FORMAT = "YYYY-MM-DD";

const TIMES_SELECTABLE: TimeConstraint[] = [
    new TimeConstraint("", "", ""),
    new TimeConstraint(moment().startOf("week").add(1, "day").format(MOMENT_FORMAT), moment().add(1, "day").format(MOMENT_FORMAT), "this week"),
    new TimeConstraint(moment().startOf("week").subtract(1, 'week').add(1, "day").format(MOMENT_FORMAT), moment().startOf("week").add(1, "day").format(MOMENT_FORMAT), "last week"),
]
