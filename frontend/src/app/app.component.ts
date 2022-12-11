import { Component } from '@angular/core';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css'],
})
export class AppComponent {
  title = 'frontend';
  show_star_window: boolean = true;

  lookup_name: string = "";

  lookup_star($event: string) {
    this.lookup_name = $event
    this.show_star_window = true;
  }

  lookup_constellation($event: string) {
    this.lookup_name = $event
    this.show_star_window = false;
  }

  clear_lookup(){
    this.lookup_name = '';
  }

  handle_event_show_star_window(){
    this.clear_lookup();
    this.show_star_window = true;
  }


  handle_event_show_constellation_window(){
    this.clear_lookup();
    this.show_star_window = false;
  }
}
