import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ApplicationNoticeComponent } from './application-notice.component';

describe('ApplicationNoticeComponent', () => {
  let component: ApplicationNoticeComponent;
  let fixture: ComponentFixture<ApplicationNoticeComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ApplicationNoticeComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(ApplicationNoticeComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
