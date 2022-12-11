import { NgTemplateOutlet } from '@angular/common';
import { Component, ElementRef, OnInit, ViewChild, Input, Output, EventEmitter } from '@angular/core';

import { Constellation } from './constellation';
import { ConstellationService } from './constellation.service';

@Component({
  selector: 'app-constellation',
  templateUrl: './constellation.component.html',
  providers: [ConstellationService],
  styleUrls: ['./constellation.component.css'],
})
export class ConstellationComponent implements OnInit {
  constellation_list: Constellation[] = [];
  selected_constellation: Constellation | undefined;
  edited_constellation: Constellation | undefined; // the constellation currently being edited
  new_constellation: Constellation | undefined;

  constellation_name_search: string = '';

  @Input() lookup_constellation_name: string = '';
  @Output() lookup_star_emitter = new EventEmitter<string>();

  constructor(private constellationService: ConstellationService) {}

  ngOnInit() {
    if(this.lookup_constellation_name === ''){
      this.get_constellation_list();
    }
    else{
      this.search_constellation_name(this.lookup_constellation_name);
    }
  }

  handle_event_select_constellation(constellation_id: number) {
    this.edited_constellation = undefined;
    this.new_constellation = undefined;
    this.selected_constellation = this.constellation_list.find(
      (constellation) => constellation.constellation_id === constellation_id
    );
  }

  handle_event_edit_constellation(){
    this.edited_constellation = this.selected_constellation;
    this.selected_constellation = undefined;
    this.new_constellation = undefined;
  }

  handle_event_confirm_edit() {
    if(this.edited_constellation && this.edited_constellation.name){
      this.update_constellation();  
      this.selected_constellation = this.edited_constellation;
    }
    else{
      this.selected_constellation = this.constellation_list[0];
    }
    this.edited_constellation = undefined;
    this.new_constellation = undefined;
  }

  handle_event_new_constellation() {
    this.selected_constellation = undefined;
    this.edited_constellation = undefined;
    this.new_constellation = { name: 'New constellation' } as Constellation;
  }

  handle_event_confirm_new() {
    if(this.new_constellation && this.new_constellation.name){
      this.add_constellation();
    }
    else{
      this.selected_constellation = this.constellation_list[0];
    }
    this.edited_constellation = undefined;
    this.new_constellation = undefined;
  }

  handle_event_delete_constellation() {
    this.delete_constellation(this.selected_constellation!);
    this.selected_constellation = this.constellation_list[0];
    this.edited_constellation = undefined;
    this.new_constellation = undefined
  }

  handle_event_lookup_star(name: string){
    // Open the information pannel about a star of the selected constellation
    if(name){
      this.lookup_star_emitter.emit(name);
    }
  }

  handle_event_cancel(){
    // Cancel the current action and reset the selected star
    if(this.lookup_constellation_name === ''){
      this.get_constellation_list();
    }
    else{
      this.search_constellation_name(this.lookup_constellation_name);     
    }
    this.edited_constellation = undefined;
    this.new_constellation = undefined;
  }

  get_constellation_list(): void {
    this.constellationService
      .get_constellation_list('')
      .subscribe(
        (constellation_list) => (
          (this.constellation_list = constellation_list), (this.selected_constellation = constellation_list[0])
        )
      );
  }

  add_constellation(): void{
    this.constellationService
          .add_constellation(this.new_constellation!)
          .subscribe((constellation) => ((this.constellation_list.push(constellation)), (this.selected_constellation = constellation)));
  }

  delete_constellation(constellation: Constellation): void {
    this.constellation_list = this.constellation_list.filter((d) => d !== constellation);
    this.constellationService.delete_constellation_by_id(constellation.constellation_id).subscribe();
  }

  search_constellation_name(constellation_name: string) {
    this.edited_constellation = undefined;
    if (constellation_name && constellation_name !== '') {
      this.constellationService
        .get_constellation_list(constellation_name)
        .subscribe((constellation_list) => (this.constellation_list = constellation_list, this.constellation_list.length > 0 ? this.selected_constellation = this.constellation_list[0] : ""));
    } else {
      this.get_constellation_list();
    }
  }

  update_constellation() {
    this.constellationService
      .update_constellation({
        ...this.edited_constellation!
      })
      .subscribe((constellation) => {
        // replace the constellation in the constellation_list list with update from server
        const index = constellation
          ? this.constellation_list.findIndex((d) => d.constellation_id === constellation.constellation_id)
          : -1;
        if (index > -1) {
          this.constellation_list[index] = constellation;         
        }
    }); 
  }
}