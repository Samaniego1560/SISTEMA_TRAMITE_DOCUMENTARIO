import { ComponentFixture, TestBed } from '@angular/core/testing';

import { RequestsReceivedComponent } from './requests-received.component';

describe('RequestsReceivedComponent', () => {
  let component: RequestsReceivedComponent;
  let fixture: ComponentFixture<RequestsReceivedComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [RequestsReceivedComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(RequestsReceivedComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
