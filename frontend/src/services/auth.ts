import config from '../config';
const authLoginGithub = async () => {
    try {
        const url = new URL('https://github.com/login/oauth/authorize');
        url.searchParams.append('client_id', config.githubClientId);
        url.searchParams.append('scope', 'user:email');
        window.location.href = url.toString();
    } catch (error) {
        console.log(error);

    }
}

export { authLoginGithub }