import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { HttpHeaders } from '@angular/common/http';

import { Observable } from 'rxjs';
import { catchError } from 'rxjs/operators';

import { Constellation } from './constellation';
import { HttpErrorHandler, HandleError } from '../http-error-handler.service';

const httpOptions = {
  headers: new HttpHeaders({
    'Content-Type': 'application/json',
    Authorization: 'my-auth-token',
  }),
};

@Injectable()
export class ConstellationService {
  constellationUrl = 'http://localhost:3000/celeste/constellation'; // URL to web api
  private handleError: HandleError;

  constructor(private http: HttpClient, httpErrorHandler: HttpErrorHandler) {
    this.handleError = httpErrorHandler.createHandleError('ConstellationService');
  }

  /** GET constellations list from the server, optional name to filter the list */
  get_constellation_list(name: string): Observable<Constellation[]> {
    name = name.trim();
    let url = `${this.constellationUrl}`;
    if (name !== '' && name !== null) {
      url = `${this.constellationUrl}?name=${name}`;
    }
    return this.http
      .get<Constellation[]>(url)
      .pipe(catchError(this.handleError('get_constellation_list', [])));
  }

  /** POST: add a new constellation to the database */
  add_constellation(constellation: Constellation): Observable<Constellation> {
    const url = `${this.constellationUrl}`;
    return this.http
      .post<Constellation>(url, constellation, httpOptions)
      .pipe(catchError(this.handleError('add_constellation', constellation)));
  }

  /** PUT: update the constellation on the server. Returns the updated constellation upon success. */
  update_constellation(constellation: Constellation): Observable<Constellation> {
    const url = `${this.constellationUrl}`;
    return this.http
      .put<Constellation>(url, constellation, httpOptions)
      .pipe(catchError(this.handleError('update_constellation', constellation)));
  }

  /** GET constellation with id from the server */
  get_constellation_by_id(constellation_id: number): Observable<Constellation> {
    const url = `${this.constellationUrl}/${constellation_id}`;
    return this.http
      .get<Constellation>(url, httpOptions)
      .pipe(catchError(this.handleError<Constellation>('get_constellation_by_id')));
  }

  /** DELETE: delete the constellation from the server */
  delete_constellation_by_id(constellation_id: number): Observable<unknown> {
    const url = `${this.constellationUrl}/${constellation_id}`;
    return this.http
      .delete(url, httpOptions)
      .pipe(catchError(this.handleError('delete_constellation_by_id')));
  }

  add_star_to_constellation(constellation_id: number, star_id: number): Observable<Constellation> {
    const url = `${this.constellationUrl}/link?star_id=${star_id}&constellation_id=${constellation_id}`;
    return this.http
      .post<Constellation>(url, httpOptions)
      .pipe(catchError(this.handleError<Constellation>('add_star_to_constellation')));
  };

  remove_star_from_constellation(constellation_id: number, star_id: number): Observable<Constellation> {
    const url = `${this.constellationUrl}/link?star_id=${star_id}&constellation_id=${constellation_id}`;
    return this.http
      .delete<Constellation>(url, httpOptions)
      .pipe(catchError(this.handleError<Constellation>('remove_star_from_constellation')));
  };
}
