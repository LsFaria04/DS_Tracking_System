import React from "react";

interface ModalProps {
  show: boolean;
  onClose: () => void;
  onUpdate: () => void;
  isUpdating: boolean;
  children: React.ReactNode;
}

export default function UpdateModal ({ show, onClose, onUpdate, isUpdating,  children }: ModalProps){
  if (!show) return null; // don't render if not visible

  return (
    <div>
      <div className="fixed inset-0 z-1000 flex items-center justify-center bg-black/75 backdrop-blur-xs">
        <div className="bg-gray-200 dark:bg-gray-800  rounded-lg shadow-lg p-6 w-[400px]">
          {children}
          <div className="flex flex-row justify-evenly gap-4">
              <div className="flex flex-row justify-evenly gap-4 mt-4">
                <button
                  onClick={onUpdate}
                  disabled={isUpdating}
                  className="w-full px-4 py-2 max-h-10 bg-blue-500 hover:bg-blue-600 
                            disabled:bg-gray-400 text-white text-sm font-medium rounded-lg 
                            transition-colors disabled:cursor-not-allowed flex items-center justify-center gap-2"
                >
                  {isUpdating ? (
                    <>
                      <span className="inline-block w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
                      Updating...
                    </>
                  ) : (
                    "Update"
                  )}
                </button>

                <button
                  onClick={onClose}
                  disabled={isUpdating}
                  className="w-full px-4 py-2 max-h-10 bg-gray-300 hover:bg-gray-400 
                            text-black text-sm font-medium rounded-lg transition-colors 
                            disabled:cursor-not-allowed flex items-center justify-center gap-2"
                >
                  Cancel
                </button>
              </div>
          </div>
         
        </div>
      </div>
    </div>
  );
}