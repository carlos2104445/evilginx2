import React from 'react';

type Props = {
  open: boolean;
  title?: string;
  message: string;
  onConfirm: () => void;
  onCancel: () => void;
};

export default function ConfirmDialog({ open, title = 'Please Confirm', message, onConfirm, onCancel }: Props) {
  if (!open) return null;
  return (
    <div style={{
      position: 'fixed', inset: 0, background: 'rgba(0,0,0,0.4)',
      display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 9998
    }}>
      <div style={{ background: '#fff', color: '#111', borderRadius: 8, width: 380, padding: 16, boxShadow: '0 6px 24px rgba(0,0,0,0.2)' }}>
        <h3 style={{ marginTop: 0 }}>{title}</h3>
        <p>{message}</p>
        <div style={{ display: 'flex', justifyContent: 'flex-end', gap: 8 }}>
          <button onClick={onCancel}>Cancel</button>
          <button onClick={onConfirm} style={{ background: '#d64545', color: '#fff', border: 'none', padding: '6px 10px', borderRadius: 4 }}>Confirm</button>
        </div>
      </div>
    </div>
  );
}
