import { Star } from '../star/star';
export interface Constellation {
    constellation_id: number;
    name: string;
    description: string;
    right_ascension: string;
    declination: string;
    area: number;
    stars: Star[];
  }