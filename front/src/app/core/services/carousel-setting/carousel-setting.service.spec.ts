import { TestBed } from '@angular/core/testing';

import { CarouselSettingService } from './carousel-setting.service';

describe('CarouselSettingService', () => {
    let service: CarouselSettingService;

    beforeEach(() => {
        TestBed.configureTestingModule({});
        service = TestBed.inject(CarouselSettingService);
    });

    it('should be created', () => {
        expect(service).toBeTruthy();
    });
});
