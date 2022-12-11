import { NgTemplateOutlet } from '@angular/common';
import { Component, ElementRef, OnInit, ViewChild, Output, EventEmitter, Input } from '@angular/core';
import { AppComponent } from '../app.component';
import { Constellation } from '../constellation/constellation';

import { Star } from './star';
import { StarService } from './star.service';
import { ConstellationService } from '../constellation/constellation.service';

@Component({
  selector: 'app-star',
  templateUrl: './star.component.html',
  providers: [StarService, ConstellationService],
  styleUrls: ['./star.component.css'],
})
export class StarComponent implements OnInit {
  star_list: Star[] = [];
  constellation_list: Constellation[] = [];
  selected_star: Star | undefined;
  edited_star: Star | undefined; // the star currently being edited
  new_star: Star | undefined;
  constellation: Constellation | undefined;

  star_name_search: string = '';

  @Input() lookup_star_name: string = '';
  @Output() lookup_constellation_emitter = new EventEmitter<string>();

  constructor(private starService: StarService, private constellationService: ConstellationService) {}

  ngOnInit() {
    if(this.lookup_star_name === ''){
      this.get_star_list();
    }
    else{
      this.search_star_name(this.lookup_star_name);     
    }
    this.get_constellation_list();
  }

  handle_event_select_star(star_id: number) {
    this.edited_star = undefined;
    this.new_star = undefined;
    this.selected_star = this.star_list.find(
      (star) => star.star_id === star_id
    );
  }

  handle_event_edit_star(){
    this.constellation = this.selected_star?.constellation;
    this.edited_star = this.selected_star;
    this.selected_star = undefined;
    this.new_star = undefined;
  }

  handle_event_confirm_edit() {
    if(this.edited_star && this.edited_star.name){
      this.update_star();  
      this.selected_star = this.edited_star;
    }
    else{
      this.selected_star = this.star_list[0];
    }
    this.constellation = undefined;
    this.edited_star = undefined;
    this.new_star = undefined;
  }

  handle_event_new_star() {
    this.selected_star = undefined;
    this.edited_star = undefined;
    this.new_star = { name: 'New star' } as Star;
  }

  handle_event_confirm_new() {
    if(this.new_star && this.new_star.name){
      this.add_star();
    }
    else{
      this.selected_star = this.star_list[0];
    }
    this.constellation = undefined;
    this.edited_star = undefined;
    this.new_star = undefined;
  }

  handle_event_delete_star() {
    this.delete_star(this.selected_star!);
    this.selected_star = this.star_list[0];
    this.edited_star = undefined;
    this.new_star = undefined;
  }

  handle_event_lookup_constellation(){
    if(this.selected_star && this.selected_star.constellation){
      this.lookup_constellation_emitter.emit(this.selected_star.constellation.name);
    }
  }

  handle_event_cancel(){
    if(this.lookup_star_name === ''){
      this.get_star_list();
    }
    else{
      this.search_star_name(this.lookup_star_name);     
    }
    this.edited_star = undefined;
    this.new_star = undefined;
  }

  get_star_list(): void {
    this.starService
      .get_star_list('')
      .subscribe(
        (star_list) => (
          (this.star_list = star_list), (this.selected_star = star_list[0])
        )
      );
  }

  get_constellation_list(): void {
    this.constellationService
      .get_constellation_list('')
      .subscribe(
        (constellation_list) => (
          (this.constellation_list = constellation_list)
        )
      );
  }

  add_star(): void{ 
    this.constellation = this.new_star!.constellation!;
    console.log(this.new_star);
    console.log(this.constellation);
    this.starService
          .add_star(this.new_star!)
          .subscribe((star) => (
            (this.link_star_to_constellation(star)),
            (this.star_list.push(star)),
            (this.selected_star = star)
            ));
  }

  link_star_to_constellation(star: Star): void{
    if(this.constellation){
      this.constellationService.add_star_to_constellation(this.constellation.constellation_id, star.star_id).subscribe();
    }
  }

  delete_star(star: Star): void {
    this.star_list = this.star_list.filter((d) => d !== star);
    this.starService.delete_star_by_id(star.star_id).subscribe();
  }

  search_star_name(star_name: string) {
    this.edited_star = undefined;
    if (star_name && star_name !== '') {
      this.starService
        .get_star_list(star_name)
        .subscribe((star_list) => (this.star_list = star_list, this.star_list.length > 0 ? this.selected_star = this.star_list[0] : ""));
    } else {
      this.get_star_list();
    }
  }

  update_star() {
    if(this.constellation && !this.edited_star?.constellation){
      this.constellationService.remove_star_from_constellation(this.constellation.constellation_id, this.edited_star!.star_id).subscribe();
    }
    else if(this.edited_star?.constellation && this.edited_star?.constellation !== this.constellation){
      this.constellationService.add_star_to_constellation(this.edited_star.constellation.constellation_id, this.edited_star.star_id).subscribe();
    }
    this.starService
      .update_star({
        ...this.edited_star!
      })
      .subscribe((star) => {
        // replace the star in the star_list list with update from server
        const index = star
          ? this.star_list.findIndex((d) => d.star_id === star.star_id)
          : -1;
        if (index > -1) {
          this.star_list[index] = star;         
        }
    }); 
  }
}
