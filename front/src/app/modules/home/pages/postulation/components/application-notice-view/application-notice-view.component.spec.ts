import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ApplicationNoticeViewComponent } from './application-notice-view.component';

describe('ApplicationNoticeViewComponent', () => {
  let component: ApplicationNoticeViewComponent;
  let fixture: ComponentFixture<ApplicationNoticeViewComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ApplicationNoticeViewComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(ApplicationNoticeViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
