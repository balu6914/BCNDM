import { HttpParams } from '@angular/common/http';

const keys = ['x0', 'y0', 'x1', 'y1', 'x2', 'y2', 'x3', 'y3'];

export class Query {
  public page: number;
  public limit: number;
  public name?: string;
  public streamType?: string;
  public owner?: string;
  public minPrice?: number;
  public maxPrice?: number;
  private _coords?: Map<string, number> | null;
  constructor(
  ) {
    // Set default values.
    this.page = 0;
    this.limit = 20;
    this._coords = new Map<string, number>();
  }

  setPoint(key: string, value: number) {
    // Create new map if not exist.
    this._coords = this._coords || new Map();
    // Set value checking keys.
    if (keys.find(k => k === key)) {
      this._coords.set(key, value);
    }
  }

  get coords(): Map<string, number> {
    return this._coords;
  }

  generateQuery(): HttpParams {
    let params = new HttpParams();
    params = params.set('page', this.page.toString());
    params = params.set('limit', this.limit.toString());

    if (this.minPrice) {
      params = params.set('minPrice', this.minPrice.toString());
    }

    if (this.maxPrice) {
      params = params.set('maxPrice', this.maxPrice.toString());
    }

    if (this.name) {
      params = params.set('name', this.name);
    }

    if (this.streamType) {
      params = params.set('type', this.streamType);
    }

    if (this.owner) {
      params = params.set('owner', this.owner);
    }

    Array.from(this._coords.entries()).forEach(value => {
      params = params.set(value[0], value[1].toString());
    });

    return params;
  }
}
