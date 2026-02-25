import { Injectable } from '@angular/core';
import { HttpClient } from "@angular/common/http";
import { Observable } from "rxjs";
import { IResponse } from "../../models/response";
import { ICarouselSetting } from "../../models/carousel-setting";

declare const window: any;

@Injectable({
    providedIn: 'root'
})
export class CarouselSettingService {
    private API_URL = window.__env.API_DBU;

    constructor(private _http: HttpClient) { }

    /**
     * Get all carousel settings (admin)
     */
    getCarouselSettings(): Observable<IResponse<ICarouselSetting[]>> {
        return this._http.get<IResponse<ICarouselSetting[]>>(`${this.API_URL}/carousel/settings`);
    }

    /**
     * Get only enabled carousel settings (public)
     */
    getPublicCarouselSettings(): Observable<IResponse<ICarouselSetting[]>> {
        return this._http.get<IResponse<ICarouselSetting[]>>(`${this.API_URL}/carousel/settings/public`);
    }

    /**
     * Create a new carousel setting
     */
    createCarouselSetting(formData: FormData): Observable<IResponse<ICarouselSetting>> {
        return this._http.post<IResponse<ICarouselSetting>>(`${this.API_URL}/carousel/settings/create`, formData);
    }

    /**
     * Update a single carousel setting
     */
    updateCarouselSetting(id: number, formData: FormData): Observable<IResponse<ICarouselSetting>> {
        return this._http.post<IResponse<ICarouselSetting>>(`${this.API_URL}/carousel/settings/update/${id}`, formData);
    }

    /**
     * Update multiple carousel settings (bulk update)
     */
    updateCarouselSettings(data: ICarouselSetting[]): Observable<IResponse<ICarouselSetting[]>> {
        return this._http.post<IResponse<ICarouselSetting[]>>(`${this.API_URL}/carousel/settings/update`, data);
    }

    /**
     * Delete a carousel setting
     */
    deleteCarouselSetting(id: number): Observable<IResponse<null>> {
        return this._http.delete<IResponse<null>>(`${this.API_URL}/carousel/settings/destroy/${id}`);
    }
}
