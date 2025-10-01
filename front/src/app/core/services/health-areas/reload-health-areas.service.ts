import {Subject} from "rxjs";
import {Injectable} from "@angular/core";

@Injectable({
  providedIn: 'root'
})
export class ReloadHealthAreasService {
  private _reloadNursing = new Subject<void>();
  private _reloadMedicine = new Subject<void>();
  private _reloadDentistry = new Subject<void>();
  public reloadNursingConsulting$ = this._reloadNursing.asObservable();
  public reloadMedicineConsulting$ = this._reloadMedicine.asObservable();
  public reloadDentistryConsulting$ = this._reloadDentistry.asObservable();

  public reloadNursingConsulting() {
    this._reloadNursing.next();
  }

  public reloadMedicineConsulting() {
    this._reloadMedicine.next();
  }

  public reloadDentistryConsulting() {
    this._reloadDentistry.next();
  }
}
