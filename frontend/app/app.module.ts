import { NgModule }      from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { FormsModule }   from '@angular/forms';
import { HttpModule  }   from '@angular/http';
import { MomentModule }  from 'angular2-moment';
import { ChartsModule } from 'ng2-charts/ng2-charts';
import { Ng2TableModule  } from 'ng2-table/ng2-table';

import { AppComponent }         from './app.component';
import { TimePickerComponent }  from './time-pick.component';
import { DevListComponent }     from './dev-list.component';
import { CompanyListComponent } from './company-list.component';

@NgModule({
  imports:      [
	  BrowserModule,
	  FormsModule,
	  MomentModule,
	  HttpModule,
	  ChartsModule,
	  Ng2TableModule
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

