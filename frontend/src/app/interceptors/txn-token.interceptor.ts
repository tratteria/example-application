import { Injectable } from '@angular/core';
import { HttpInterceptor, HttpRequest, HttpHandler, HttpEvent, HttpResponse } from '@angular/common/http';
import { Observable } from 'rxjs';
import { tap } from 'rxjs/operators';
import { TxnTokenService } from '../services/txn-token.service';

@Injectable()
export class TxnTokenInterceptor implements HttpInterceptor {
  constructor(private txnTokenService: TxnTokenService) {}

  intercept(request: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    return next.handle(request).pipe(
      tap(event => {
        if (event instanceof HttpResponse) {
          const txnToken = event.headers.get('Txn-Token');
          if (txnToken) {
            this.txnTokenService.setTxnToken(txnToken);
          }
        }
      })
    );
  }
}