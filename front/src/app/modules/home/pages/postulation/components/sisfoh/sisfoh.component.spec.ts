import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SisfohComponent } from './sisfoh.component';

describe('SisfohComponent', () => {
  let component: SisfohComponent;
  let fixture: ComponentFixture<SisfohComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [SisfohComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(SisfohComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
