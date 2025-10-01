import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ReleaseSettingComponent } from './release-setting.component';

describe('ReleaseSettingComponent', () => {
  let component: ReleaseSettingComponent;
  let fixture: ComponentFixture<ReleaseSettingComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ReleaseSettingComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(ReleaseSettingComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
