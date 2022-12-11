import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { HttpClientModule } from '@angular/common/http';
import { FormsModule } from '@angular/forms';

import { AppRoutingModule } from './app-routing.module';
import { HttpErrorHandler } from './http-error-handler.service';
import { MessageService } from './message.service';
import { AppComponent } from './app.component';
import { StarComponent } from './star/star.component';
import { ConstellationComponent } from './constellation/constellation.component';

@NgModule({
  declarations: [
    AppComponent,
    StarComponent,
    ConstellationComponent,
  ],
  imports: [
    BrowserModule,
    HttpClientModule,
    AppRoutingModule,
    FormsModule
  ],
  providers: [
    HttpErrorHandler,
    MessageService
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
