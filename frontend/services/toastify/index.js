import { toast } from "react-toastify";

const configuration = {
  position: "top-right",
  autoClose: 5000,
  hideProgressBar: false,
  closeOnClick: true,
  pauseOnHover: true,
  draggable: true
};

export const showErrorMessage = message => toast.error(message, configuration);

export const showSuccessMessage = message =>
  toast.success(message, configuration);

export const showInfoMessage = message => toast.info(message, configuration);
