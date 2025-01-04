import { API_URL } from "./App"
import {useNavigate, useParams} from "react-router-dom"

const ConfirmationPage = () => {
    const {token = ""} = useParams()
    const redirect = useNavigate()
    console.log(token)
    const handleConfirm = async () => {
        const response = await fetch(`${API_URL}/users/activate/${token}`, {
            method:"PUT",
        
        })
        if (response.ok){
            redirect("/")
        }else {
            alert("Failed to confirm user")
        }
    }
  return (
    <div>
        <h1>Confirmation</h1>
        <button onClick={handleConfirm}>Click to confirm</button>
    </div>
  )
}

export default ConfirmationPage