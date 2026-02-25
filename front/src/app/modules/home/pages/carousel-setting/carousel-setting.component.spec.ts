import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CarouselSettingComponent } from './carousel-setting.component';

describe('CarouselSettingComponent', () => {
    let component: CarouselSettingComponent;
    let fixture: ComponentFixture<CarouselSettingComponent>;

    beforeEach(async () => {
        await TestBed.configureTestingModule({
            imports: [CarouselSettingComponent]
        })
            .compileComponents();

        fixture = TestBed.createComponent(CarouselSettingComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
