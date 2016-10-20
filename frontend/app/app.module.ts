import { NgModule }      from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { FormsModule }   from '@angular/forms';
import { HttpModule  }   from '@angular/http';
import { MomentModule }  from 'angular2-moment';

import { AppComponent }         from './app.component';
import { TimePickerComponent }  from './time-pick.component';
import { DevListComponent }     from './dev-list.component';
import { CompanyListComponent } from './company-list.component';

@NgModule({
  imports:      [
	  BrowserModule,
	  FormsModule,
	  MomentModule,
	  HttpModule
  ],
  declarations: [
	  AppComponent,
	  TimePickerComponent,
	  DevListComponent,
	  CompanyListComponent
  ],
  bootstrap: [ AppComponent ]
})

export class AppModule { }

