import { Constellation } from '../constellation/constellation';
export interface Star {
    star_id: number;
    name: string;
    description: string;
    right_ascension: string;
    declination: string;
    apparent_magnitude: number;
    mass: number;
    radius: number;
    age: number;
    constellation: Constellation;
  }