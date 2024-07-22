import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class TxnTokenService {
  private txnTokenSubject = new BehaviorSubject<string | null>(null);
  txnToken$ = this.txnTokenSubject.asObservable();

  setTxnToken(token: string) {
    this.txnTokenSubject.next(token);
  }

  clearTxnToken() {
    this.txnTokenSubject.next(null);
  }
}