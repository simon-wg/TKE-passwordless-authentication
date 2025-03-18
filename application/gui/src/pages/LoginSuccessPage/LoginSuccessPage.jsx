// import React, { useEffect, useState } from "react";
// import { useNavigate } from "react-router-dom";
// import useFetchUser from "../../hooks/useFetchUser";
// import NotesApp from "../../components/NotesApp";

// const LoginSuccessPage = () => {
//   const user = useFetchUser();

//   return (
//     <div>
//       <NotesApp />
//     </div>
//   );
// };

// export default LoginSuccessPage;

import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import useFetchUser from "../../hooks/useFetchUser";
import NotesApp from "../../components/NotesApp";

const LoginSuccessPage = () => {
  const user = useFetchUser();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!loading && user === null) {
      navigate("/register");
    }
  }, [user, navigate]);

  if (user == null) {
    return <p>
      Loading...
    </p>
  }

  return (
    <div>
      <NotesApp />
    </div>
  );
};

export default LoginSuccessPage;