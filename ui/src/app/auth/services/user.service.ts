// Imports
import { Injectable }     from '@angular/core';
import { Http, Response, Headers, RequestOptions } from '@angular/http';
import { environment } from '../../../environments/environment';
import { User } from '../../common/interfaces/user.interface'
import {Observable} from 'rxjs/Rx';

// Import RxJs required methods
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/catch';

@Injectable()
export class UserService {
     // Resolve HTTP using the constructor
     constructor (private http: Http) {}

     // Fetch  user
     getUser() : Observable<User[]> {
         return this.http.get(`${environment.API_AUTH}/users`)
                         .map((res:Response) => res.json())
                         .catch((error:any) => Observable.throw(error.json().error || 'Server error'));
     }
     addUser (data): Observable<User[]> {
        return this.http.post(`${environment.API_AUTH}/users`, data)
                         .map((res:Response) => res)
                         .catch((error:any) => Observable.throw(error.json().error || 'Server error'));
    }
     removeUser (id:string): Observable<User[]> {
        return this.http.delete(`${environment.API_AUTH}/users`)
                         .map((res:Response) => res.json())
                         .catch((error:any) => Observable.throw(error.json().error || 'Server error'));
    }
     updateUser (id:string, data): Observable<User[]> {
        return this.http.put(`${environment.API_AUTH}/users`, data)
                         .map((res:Response) => res.json())
                         .catch((error:any) => Observable.throw(error.json().error || 'Server error'));
    }

}
