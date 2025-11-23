import React from "react";

interface ModalProps {
  show: boolean;
  onClose: () => void;
  onUpdate: () => void;
  children: React.ReactNode;
}

export default function UpdateModal ({ show, onClose, onUpdate, children }: ModalProps){
  if (!show) return null; // don't render if not visible

  return (
    <div>
      <div >
        {children}
        <button onClick={onUpdate} >Update</button>
        <button onClick={onClose} >Cancel</button>
      </div>
    </div>
  );
}