const fs = require('fs');
const https = require('https');
const yargs = require('yargs');

const argv = yargs(process.argv.slice(2))
    .scriptName("qb-tracker-updater")
    .usage('$0 <qBittorrent config file path> [profile]')
    .option('profile', {
        description: 'Select which tracker list update with',
        choices: ['best', 'best_ip', 'all', 'all_http', 'all_https', 'all_ws'],
        default: 'best_ip',
    })
    .help()
    .argv

const options = {
    hostname: 'raw.githubusercontent.com',
    path: '',
    headers: {
        'User-Agent': 'Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36'
    },
    method: 'GET',
    port: 443,
}

if (argv.help) {
    console.log('tracker_updater <qBittorrent config file path> [profile]')
    process.exit(0)
}

switch (argv.profile) {
    case 'best':
        options.path = '/ngosang/trackerslist/master/trackers_best.txt';
        break;
    case 'best_ip':
        options.path = '/ngosang/trackerslist/master/trackers_best_ip.txt';
        break;
    case 'all':
        options.path = '/ngosang/trackerslist/master/trackers_all_udp.txt';
        break;
    case 'all_http':
        options.path = '/ngosang/trackerslist/master/trackers_all_http.txt';
        break;
    case 'all_https':
        options.path = '/ngosang/trackerslist/master/trackers_all_https.txt';
        break;
    case 'all_ws':
        options.path = '/ngosang/trackerslist/master/trackers_all_ws.txt';
        break;
    case 'best_ip':
        options.path = '/ngosang/trackerslist/master/trackers_best_ip.txt';
        break;
    case 'all_ip':
        options.path = '/ngosang/trackerslist/master/trackers_best_ip.txt';
        break;
    default:
        options.path = '/ngosang/trackerslist/master/trackers_best_ip.txt';
        break;
}


if (process.argv.length < 3) {
    console.log()
    console.error('You must specify the qBittorrent config file path')
    console.log('  e.g. node tracker_updater.js /home/user/.config/qBittorrent/qBittorrent.conf [profile]\n')
    console.log('Run with --help to show more info\n')
    process.exit(-1)
}

const configFilePath = process.argv[2];

function getTrackerList() {
    return new Promise((resolve, reject) => {
        let chunks = [];
        https.get(options, res => {
            res.on('data', chunk => {
                chunks.push(chunk)
            });
            res.on('end', () => {
                const buffer = Buffer.concat(chunks)
                resolve(buffer.toString())
            });
            res.on('error', () => reject());
        });
    });
}

function trackerListToQBittorrent(trackers) {
    return trackers
        .split('\n')
        .filter(line => line !== '')
        .reduce((prev, curr) => `${prev}\\n${curr}`)
}

function readAndUpdateConfigFile() {
    return new Promise((resolve, reject) => {
        fs.readFile(configFilePath, 'utf-8', (err, file) => {
            if (err) {
                console.error(configFilePath, 'File not found')
                reject()
                return
            }
            getTrackerList().then(trackers => {
                const updated = file.split('\n').map(line => {
                    if (line.includes('Bittorrent\\TrackersList')) {
                        const [prop, _] = line.split('=')
                        return `${prop}=${trackerListToQBittorrent(trackers)}`
                    }
                    return line
                })
                resolve(updated.reduce((prev, curr) => `${prev}\n${curr}`))
            })
        });
    });
}

readAndUpdateConfigFile()
    .then(content => {
        fs.writeFile(configFilePath, content, err => {
            if (err) {
                console.error(err);
            }
        });
        console.log(`Updated ${configFilePath} with profile ${argv.profile}`)
    })
    .catch(e => console.error('ERROR: Cannot write config file!'))