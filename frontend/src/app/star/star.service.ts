import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { HttpHeaders } from '@angular/common/http';

import { Observable } from 'rxjs';
import { catchError } from 'rxjs/operators';

import { Star } from './star';
import { HttpErrorHandler, HandleError } from '../http-error-handler.service';

const httpOptions = {
  headers: new HttpHeaders({
    'Content-Type': 'application/json',
    Authorization: 'my-auth-token',
  }),
};

@Injectable()
export class StarService {
  starUrl = 'http://localhost:3000/celeste/star'; // URL to web api
  private handleError: HandleError;

  constructor(private http: HttpClient, httpErrorHandler: HttpErrorHandler) {
    this.handleError = httpErrorHandler.createHandleError('StarService');
  }

  /** GET stars list from the server */
  get_star_list(name: string): Observable<Star[]> {
    name = name.trim();
    let url = `${this.starUrl}`;
    if (name !== '' && name !== null) {
      url = `${this.starUrl}?name=${name}`;
    }
    return this.http
      .get<Star[]>(url)
      .pipe(catchError(this.handleError('get_star_list', [])));
  }

  /** POST: add a new star to the database */
  add_star(star: Star): Observable<Star> {
    const url = `${this.starUrl}`;
    return this.http
      .post<Star>(url, star, httpOptions)
      .pipe(catchError(this.handleError('add_star', star)));
  }

  /** PUT: update the star on the server. Returns the updated star upon success. */
  update_star(star: Star): Observable<Star> {
    const url = `${this.starUrl}`;
    return this.http
      .put<Star>(url, star, httpOptions)
      .pipe(catchError(this.handleError('update_star', star)));
  }

  /** GET star with id from the server */
  get_star_by_id(star_id: number): Observable<Star> {
    const url = `${this.starUrl}/${star_id}`;
    return this.http
      .get<Star>(url, httpOptions)
      .pipe(catchError(this.handleError<Star>('get_star_by_id')));
  }

  /** DELETE: delete the star from the server */
  delete_star_by_id(star_id: number): Observable<unknown> {
    const url = `${this.starUrl}/${star_id}`;
    return this.http
      .delete(url, httpOptions)
      .pipe(catchError(this.handleError('delete_star_by_id')));
  }
}
