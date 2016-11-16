import { Component, OnInit } from '@angular/core';

import * as moment from 'moment/moment';
import 'chartjs';

import { Developer, TimeConstraint } from './models'
import { CompanyService } from './company.service'
import { SettingsService } from './settings.service'

@Component({
  selector: 'my-app',
  template: `
    <h1>STATKUBE</h1>
	<time-picker [selected_time]="selected_time"></time-picker>
	<button type="button" class="btn btn-primary" (click)="filter()">Filter by date</button>
	<company-list [comps]="comps"></company-list>
	<!--dev-list [devs]="devs"></dev-list-->
  `,
  providers: [CompanyService, SettingsService]
})

export class AppComponent implements OnInit{
  comps: Developer[];
  devs = DEVS;
  selected_time = new TimeConstraint("2016-05-01", moment().add(1, "day").format("YYYY-MM-DD"), "");

  constructor(private companyService: CompanyService, private settingsService: SettingsService) {};

  getCompanies(): void {
	  this.companyService.getCompanies(this.selected_time, this.settingsService.getSettings()).then(comps => this.comps = comps);
  }

  filter(): void{
	  this.getCompanies();
  }

  ngOnInit(): void {
	  this.getCompanies();
  }
}

const DEVS: Developer[] = [
    new Developer("super dev", 12),
    new Developer("super dev2", 13),
    new Developer("super dev3", 5),
    new Developer("super dev4", 120),
];
