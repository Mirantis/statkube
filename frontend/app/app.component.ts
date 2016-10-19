import { Component } from '@angular/core';
import * as moment from 'moment/moment';
import { Developer, TimeConstraint } from './models'

@Component({
  selector: 'my-app',
  template: `
    <h1>STATKUBE</h1>
	<time-picker [selected_time]="selected_time"></time-picker>
	<company-list [comps]="comps"></company-list>
	<dev-list [devs]="devs"></dev-list>
  `
})

export class AppComponent {
  comps = COMPANIES;
  devs = DEVS;
  selected_time = new TimeConstraint("", "", "");
}

const DEVS: Developer[] = [
    new Developer("super dev", 12),
    new Developer("super dev2", 13),
    new Developer("super dev3", 5),
    new Developer("super dev4", 120),
];

const COMPANIES: Developer[] =[
    new Developer("Manimuru", 45561),
    new Developer("CoreOS", 2),
    new Developer("GlueGL", 444),
]

